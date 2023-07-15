import { Content } from "../../../api/domain";

interface PropsType {
    content: Content;
}
const SourcePanel: React.FC<PropsType> = ({ content }) => {
    let source = "";
    for (const headerKey in content.Headers) {
        const headerValue = content.Headers[headerKey];
        source += headerKey + ": " + headerValue + "\n";
    }
    source += "\n";
    source += content.Body;
    return <pre style={{ paddingBottom: 0, overflow: "auto" }}>{source}</pre>;
};

export default SourcePanel;
