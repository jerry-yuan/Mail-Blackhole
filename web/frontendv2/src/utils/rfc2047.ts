import { toByteArray as base64Decode } from "base64-js";
const RFC2047_PATTERN = /=\?([\w_-]+)\?([BbQq])\?([^?]+)\?=/;
const RFC2047_SPLITER = /(\?==\?)/g;

export const isRFC2047 = (str: string) => {
    return RFC2047_PATTERN.test(str);
};

const qEncodingDecode = (encodedStr: string): ArrayBuffer => {
    // 将=XX形式的编码转换为二进制数据
    const byteArray = new Uint8Array(encodedStr.length);
    let j = 0;
    for (let i = 0; i < encodedStr.length; i++) {
        // 遇到=，读取下两个字符并转换为10进制数值
        if (encodedStr[i] === "=" && i + 2 < encodedStr.length) {
            const byteVal = parseInt(encodedStr.substring(i + 1, i + 3), 16);
            byteArray[j++] = byteVal;
            i += 2;
        } else {
            // 直接将字符转换为ASCII码并存储到数组中
            byteArray[j++] = encodedStr.charCodeAt(i);
        }
    }
    // 将字节数组转换为ArrayBuffer对象并返回
    return byteArray.buffer;
};

export const decode = (encodedStr: string) => {
    const rfc2047Strs = encodedStr.split(RFC2047_SPLITER);

    let decodedStr = "";
    for (const rfc2047Str of rfc2047Strs) {
        const matches = rfc2047Str.match(RFC2047_PATTERN);
        if (!matches) {
            continue;
        }
        // split RFC2047 string
        const charset = matches[1];
        const binaryEncodeMode = matches[2];
        const encodedText = matches[3];

        // binary decode
        let encodedBinary: ArrayBuffer;
        if (binaryEncodeMode.toUpperCase() === "Q") {
            // Q-Encoding
            encodedBinary = qEncodingDecode(encodedText);
        } else if (binaryEncodeMode.toUpperCase() === "B") {
            // Base64
            encodedBinary = base64Decode(encodedText);
        } else {
            throw new Error("Unsupported binary encode mode:" + binaryEncodeMode);
        }

        // convert into string
        decodedStr += new TextDecoder(charset).decode(encodedBinary);
    }

    if (decodedStr === "") {
        decodedStr = encodedStr;
    }

    return decodedStr;
};
