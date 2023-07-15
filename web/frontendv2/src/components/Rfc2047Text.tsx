import { decode as rfc2047Decode } from "../utils/rfc2047";
interface PropsType {
    children: string | string[];
}

const Rfc2047Text: React.FC<PropsType> = function ({ children }) {
    let decodedStr = "";
    let content = "";
    if (typeof children === "string") {
        content = children;
    } else if (typeof children === "object") {
        content = children.join("");
    }
    while (content.length > 0) {
        const beginPatternIndex = content.indexOf("=?");
        const endPatternIndex = content.indexOf("?=");

        if (beginPatternIndex >= 0 && endPatternIndex >= 0) {
            // found a RFC2047 string
            decodedStr += content.substring(0, beginPatternIndex).trim();
            decodedStr += rfc2047Decode(content.substring(beginPatternIndex, endPatternIndex + 2));
            content = content.substring(endPatternIndex + 2);
        } else {
            // no RFC2047 string found
            decodedStr += content;
            content = "";
        }
    }
    return <>{decodedStr}</>;
};

export default Rfc2047Text;
