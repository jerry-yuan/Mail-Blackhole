import { Tab } from "semantic-ui-react";
import { Content, Message } from "../../api/domain";
import MimePanel from "./panels/MimePanel";
import HtmlPanel from "./panels/HtmlPanel";
import PlainPanel from "./panels/PlainPanel";
import SourcePanel from "./panels/SourcePanel";

interface PropsType {
    message: Message;
}

const hasHtml = (message: Message): boolean => {
    return hasHtmlParts([message.Content]) || (message.MIME !== null ? hasHtmlParts(message.MIME.Parts) : false);
};

const hasHtmlParts = (mimeParts: Content[]): boolean => {
    for (let mimePartIndex = 0; mimePartIndex < mimeParts.length; mimePartIndex++) {
        const mimePart = mimeParts[mimePartIndex];
        for (let headerKey in mimePart.Headers) {
            if (headerKey.toLowerCase() === "Content-Type".toLowerCase()) {
                if (mimePart.Headers[headerKey].join().match(/text\/html/i)) {
                    return true;
                }
                break;
            }
        }
    }

    return false;
};

/*
 * Panels
 */
const findFirstPart = (mimeParts: Content[], mimeType: string): Content | null => {
    for (let i = 0; i < mimeParts.length; i++) {
        const mimePart = mimeParts[i];
        if (mimePart.Headers["Content-Type"].join().startsWith(mimeType)) {
            return mimePart;
        }
        if (mimePart.MIME) {
            const firstPart = findFirstPart(mimePart.MIME.Parts, mimeType);
            if (firstPart) {
                return firstPart;
            }
        }
    }
    return null;
};

const MessageContentViewer: React.FC<PropsType> = function (params) {
    let panes = [];
    const message = params.message;
    const mimeParts = message.MIME?.Parts ?? [message.Content];
    // find first HTML Mime part
    const firstHtmlPart = findFirstPart(mimeParts, "text/html");
    // find first PlainText Mime part
    const firstPlainPart = firstHtmlPart ?? findFirstPart(mimeParts, "text/plain") ?? message.Content;
    // show html only message has html parts
    if (firstHtmlPart) {
        panes.push({ menuItem: "HTML", render: () => <HtmlPanel mimePart={firstHtmlPart} /> });
    }
    panes.push({ menuItem: "Plain", render: () => <PlainPanel mimePart={firstPlainPart} /> });
    panes.push({ menuItem: "Source", render: () => <SourcePanel content={message.Content} /> });
    panes.push({ menuItem: "MIME", render: () => <MimePanel mimeParts={params.message.MIME?.Parts} /> });
    return <Tab panes={panes} />;
};

export default MessageContentViewer;
