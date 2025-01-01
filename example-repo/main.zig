

const std = @import("std.zig");
pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    for( 0..10) |i| {
        try stdout.print("{d}\n", .{i});
    }
}