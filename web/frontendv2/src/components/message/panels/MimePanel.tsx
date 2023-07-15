import prettyBytes from "pretty-bytes";
import { Message as MessageBox, List, Icon } from "semantic-ui-react";
import { Content } from "../../../api/domain";
import { useState } from "react";
import { Collapse } from "react-collapse";
import Rfc2047Text from "../../Rfc2047Text";
import { saveAs } from "file-saver";
import { decode } from "../../../utils/rfc2047";
import contentTypeParser from "../../../utils/contentTypeParser";

interface MimePanelPropsType {
    mimeParts: Content[] | undefined;
}

interface MimePartPropsType {
    mimePart: Content;
}

interface MimePartsPropsType {
    mimeParts: Content[];
}

const MimeParts: React.FC<MimePartsPropsType> = ({ mimeParts }) => {
    return (
        <>
            {mimeParts.map((mimePart, index) => (
                <MimePart key={index} mimePart={mimePart} />
            ))}
        </>
    );
};

const MimePart: React.FC<MimePartPropsType> = ({ mimePart }) => {
    const [subCollapsed, setSubCollapsed] = useState<boolean>(true);
    const hasSubList = mimePart.MIME?.Parts !== undefined;
    const contentType = contentTypeParser.parse(mimePart.Headers["Content-Type"]?.join(""));
    // generate file icon
    let fileIcon;
    if (contentType.mime.startsWith("application/vnd.ms-excel")) {
        fileIcon = <List.Icon name="file excel outline" size="large" verticalAlign="middle" />;
    } else if (contentType.mime.startsWith("application/vnd.ms-powerpoint")) {
        fileIcon = <List.Icon name="file powerpoint outline" size="large" verticalAlign="middle" />;
    } else if (contentType.mime.startsWith("application/vnd.ms-word")) {
        fileIcon = <List.Icon name="file word outline" size="large" verticalAlign="middle" />;
    } else if (contentType.mime.startsWith("audio")) {
        fileIcon = <List.Icon name="file audio outline" size="large" verticalAlign="middle" />;
    } else if (contentType.mime.startsWith("video")) {
        fileIcon = <List.Icon name="file video outline" size="large" verticalAlign="middle" />;
    } else if (contentType.mime.startsWith("image")) {
        fileIcon = <List.Icon name="file image outline" size="large" verticalAlign="middle" />;
    } else {
        fileIcon = <List.Icon name="file outline" size="large" verticalAlign="middle" />;
    }

    // save file handler
    const downloadPart = () => {
        const fileName = decode(contentType.props.name ?? "unnamedPart.bin");
        saveAs(new Blob([mimePart.Body], { type: contentType.mime }), fileName);
    };

    return (
        <List.Item>
            {hasSubList ? (
                <List.Icon
                    name={subCollapsed ? "angle right" : "angle down"}
                    size="large"
                    verticalAlign="top"
                    onClick={() => setSubCollapsed(!subCollapsed)}
                    style={{ cursor: "pointer" }}
                />
            ) : (
                fileIcon
            )}
            <List.Content>
                <div style={{ display: "flex", justifyContent: "space-between" }}>
                    <List.Header>
                        <Rfc2047Text>{contentType.props.name ?? (hasSubList ? "Multipart partition" : "Unnamed MIME Part")}</Rfc2047Text>
                    </List.Header>
                </div>
                <List.Description style={{ display: "flex", justifyContent: "space-between" }}>
                    <div>{prettyBytes(mimePart.Size)}</div>
                    <div>
                        <Icon size="tiny" name="save" title="Save MIME part as file" onClick={downloadPart} style={{ cursor: "pointer" }}></Icon>
                    </div>
                </List.Description>
                {mimePart.MIME?.Parts !== undefined ? (
                    <Collapse isOpened={!subCollapsed}>
                        <List.List style={{ borderTopColor: "rgba(34, 36, 38, 0.15)", borderTopStyle: "solid", borderTopWidth: "0.8px" }}>
                            <MimeParts mimeParts={mimePart.MIME.Parts} />
                        </List.List>
                    </Collapse>
                ) : (
                    ""
                )}
            </List.Content>
        </List.Item>
    );
};

const MimePanel: React.FC<MimePanelPropsType> = ({ mimeParts }) => {
    return mimeParts === undefined ? (
        <MessageBox info content="No MIME Parts found in this mail." />
    ) : (
        <div style={{ paddingTop: "1em" }}>
            <List divided relaxed>
                <MimeParts mimeParts={mimeParts} />
            </List>
        </div>
    );
};
export default MimePanel;
