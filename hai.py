from collections import Counter

# Gabungkan semua angka dari 0 sampai 719 ke dalam satu string
digits = ''.join(str(i) for i in range(720))

# Hitung frekuensi tiap digit
counter = Counter(digits)

# Cetak hasil
for d in range(10):
    print(f"Digit {d}: {counter[str(d)]} kali")
