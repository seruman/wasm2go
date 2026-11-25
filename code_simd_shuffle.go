package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// A shuffle's lane selectors are instruction immediates. Specialize once per
// distinct mask, rather than interpreting selectors every time the code runs.
func (t *translator) shuffleFunc(lanes [16]byte) *ast.Ident {
	if name := t.shuffles[lanes]; name != nil {
		return name
	}
	if t.shuffles == nil {
		t.shuffles = make(map[[16]byte]*ast.Ident)
	}
	name := newID("shuffle" + strconv.Itoa(len(t.shuffles)))
	t.shuffles[lanes] = name
	t.helpers.add("v128")
	x, y := newID("x"), newID("y")
	var result ast.Expr
	switch {
	case identityShuffle(lanes, 0):
		result = x
	case identityShuffle(lanes, 16):
		result = y
	case *simdarch != "":
		result = t.shuffleLookup(name.Name, lanes)
	default:
		result = t.shuffleLanes(lanes)
	}
	t.out.Decls = append(t.out.Decls, &ast.FuncDecl{
		Name: name,
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{x, y}, Type: newID("v128")}}},
			Results: &ast.FieldList{List: []*ast.Field{{Type: newID("v128")}}},
		},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Results: []ast.Expr{result}}}},
	})
	return name
}

func identityShuffle(lanes [16]byte, base byte) bool {
	for i, lane := range lanes {
		if lane != base+byte(i) {
			return false
		}
	}
	return true
}

// Each table lookup contributes only the bytes selected from that input.
// 0x80 zeroes a lane on all targets (sign bit on x86, out of range elsewhere).
func (t *translator) shuffleLookup(name string, lanes [16]byte) ast.Expr {
	t.helpers.add("i8x16_lookup")
	var result ast.Expr
	for n, src := range []string{"x", "y"} {
		mask := &ast.CompositeLit{Type: shuffleArrayType(16, "byte")}
		used := false
		for _, lane := range lanes {
			index := 0x80
			if int(lane)/16 == n {
				index = int(lane) % 16
				used = true
			}
			mask.Elts = append(mask.Elts, shuffleIndex(index))
		}
		if !used {
			continue
		}
		id := newID(name + "_" + src)
		t.out.Decls = append(t.out.Decls, &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{
			&ast.ValueSpec{Names: []*ast.Ident{id}, Values: []ast.Expr{mask}},
		}})
		lookup := &ast.CallExpr{Fun: newID("i8x16_lookup"), Args: []ast.Expr{newID(src), &ast.UnaryExpr{Op: token.AND, X: id}}}
		if result == nil {
			result = lookup
		} else {
			result = shuffleMethod(result, "Or", lookup)
		}
	}
	return result
}

// Go 1.27 has no common byte-table lookup across all three targets. In a
// portable file, gather constant-index elements at the widest possible lane
// size (e.g. a shuffle swapping the two 64-bit halves needs just two lanes).
func (t *translator) shuffleLanes(lanes [16]byte) ast.Expr {
	width := 8
	for i, lane := range lanes {
		for width > 1 && (int(lane)%width != i%width || i%width != 0 && lane != lanes[i-1]+1) {
			width /= 2
		}
	}
	size := fmt.Sprintf("Uint%dx%d", 8*width, 16/width)
	t.packages.add("simd/archsimd")
	values := &ast.CompositeLit{Type: shuffleArrayType(16/width, "uint"+strconv.Itoa(8*width))}
	splat := true
	for i := 0; i < len(lanes); i += width {
		splat = splat && lanes[i] == lanes[0]
		lane := int(lanes[i])
		var src ast.Expr = newID("x")
		if lane >= 16 {
			src = newID("y")
		}
		if width > 1 {
			src = shuffleMethod(src, fmt.Sprintf("ReshapeToUint%ds", width*8))
		}
		values.Elts = append(values.Elts, shuffleMethod(src, "GetElem", shuffleIndex(lane%16/width)))
	}
	var result ast.Expr
	if splat {
		result = shuffleMethod(newID("archsimd"), "Broadcast"+size, values.Elts[0])
	} else {
		result = shuffleMethod(newID("archsimd"), "Load"+size+"Array", &ast.UnaryExpr{Op: token.AND, X: values})
	}
	if width > 1 {
		result = shuffleMethod(result, "ReshapeToUint8s")
	}
	return result
}

func shuffleArrayType(n int, elt string) *ast.ArrayType {
	return &ast.ArrayType{Len: shuffleIndex(n), Elt: newID(elt)}
}

func shuffleIndex(i int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(i)}
}

func shuffleMethod(x ast.Expr, method string, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: &ast.SelectorExpr{X: x, Sel: newID(method)}, Args: args}
}
