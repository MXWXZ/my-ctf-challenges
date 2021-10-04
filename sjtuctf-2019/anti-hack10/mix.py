def print_hex(bytes):
    l = [hex(int(i)) for i in bytes]
    print(" ".join(l))


pos = [0x400, 0x8000000, 0x18000000, 0x28000000, 0x38000000, 0x48000000]

with open('disk.img', 'rb+') as f:
    for i in pos:
        f.seek(i + 0x38)
        f.write(bytearray([0xef, 0x53]))
        if i != 0x400:
            f.seek(i + 0x3FC)
            f.write(bytearray([0x00, 0x00, 0x00, 0x00]))
