function encode_decode(strs) {
    // Encode
    let encoded = '';
    for (const s of strs) {
        encoded += s.length + '#' + s;
    }
    // Decode
    const res = [];
    let i = 0;
    while (i < encoded.length) {
        let j = i;
        while (encoded[j] !== '#') j++;
        const length = parseInt(encoded.slice(i, j));
        res.push(encoded.slice(j + 1, j + 1 + length));
        i = j + 1 + length;
    }
    return res;
}
