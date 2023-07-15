const decode = (str: string): Uint8Array => {
    let decodedBytes = [];
    for (let i = 0; i < str.length; ) {
        if (str.charAt(i) === "=") {
            decodedBytes.push(parseInt(str.substr(i + 1, 2), 16));
            i += 3;
        } else {
            decodedBytes.push(str.charCodeAt(i));
            ++i;
        }
    }
    return new Uint8Array(decodedBytes);
};

export { decode };
