#![no_std]

use core::{panic::PanicInfo, slice};

static mut INPUT: [u8; 65536] = [0; 65536];

#[panic_handler]
fn panic(_: &PanicInfo) -> ! {
    core::arch::wasm32::unreachable()
}

#[no_mangle]
pub extern "C" fn input() -> *mut u8 {
    core::ptr::addr_of_mut!(INPUT).cast()
}

#[no_mangle]
pub extern "C" fn find(needle: u8, n: usize) -> i32 {
    assert!(n <= 65536);
    // The caller fills INPUT before each search and doesn't mutate it during the call.
    let s = unsafe { slice::from_raw_parts(input(), n) };
    memchr::memchr(needle, s).map_or(-1, |i| i as i32)
}
