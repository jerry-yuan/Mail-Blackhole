import { toByteArray as base64Decode } from "base64-js";
import { decode as quotedPrintableDecode } from "../../../utils/quotedPrintable";
import { Content } from "../../../api/domain";
import contentTypeParser from "../../../utils/contentTypeParser";
import { Icon, Message } from "semantic-ui-react";

interface PropsType {
    mimePart: Content;
}
const PlainPanel: React.FC<PropsType> = ({ mimePart }) => {
    const contentType = contentTypeParser.parse(mimePart.Headers["Content-Type"]?.join(""));
    const contentTransferEncoding = mimePart.Headers["Content-Transfer-Encoding"]?.join("");

    let contentBinary: Uint8Array;
    // transfer decode
    let transferCodeValid = true;
    if (contentTransferEncoding) {
        if (contentTransferEncoding.toLowerCase() === "quoted-printable") {
            contentBinary = quotedPrintableDecode(mimePart.Body.replace(/=[\r\n]+/gm, ""));
        } else if (contentTransferEncoding.toLowerCase() === "base64") {
            contentBinary = base64Decode(mimePart.Body.replace(/\r?\n|\r/gm, ""));
        } else if (contentTransferEncoding in ["7bit", "8bit"]) {
            contentBinary = new TextEncoder().encode(mimePart.Body);
        } else {
            contentBinary = new TextEncoder().encode(mimePart.Body);
            transferCodeValid = false;
        }
    } else {
        contentBinary = new TextEncoder().encode(mimePart.Body);
    }
    // display transcode
    let content = new TextDecoder(contentType.props.charset ?? "utf-8").decode(contentBinary);

    return (
        <div style={{ paddingTop: "1em" }}>
            {transferCodeValid ? (
                ""
            ) : (
                <Message color="yellow">
                    <Icon name="exclamation triangle" />
                    Unsupported <code>Content-Transfer-Encoding</code> value "{contentTransferEncoding}", body is displaying without decoding.
                </Message>
            )}
            <pre style={{ overflowX: "auto", margin: 0 }}>{content}</pre>
        </div>
    );
};
export default PlainPanel;
