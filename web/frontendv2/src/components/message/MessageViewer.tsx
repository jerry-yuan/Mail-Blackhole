import { Card } from "semantic-ui-react";
import { Message } from "../../api/domain";
import MessageHeaderViewer from "./MessageHeaderViewer";
import MessageContentViewer from "./MessageContentViewer";
import Rfc2047Text from "../Rfc2047Text";
interface PropsType {
    message: Message;
}
const MessageViewer: React.FC<PropsType> = function ({ message }) {
    return (
        <Card fluid>
            <Card.Content>
                <Card.Header>
                    <Rfc2047Text>{message.Content.Headers["Subject"][0]}</Rfc2047Text>
                </Card.Header>
                <Card.Meta>
                    {message.From?.Mailbox}@{message.From?.Domain}
                </Card.Meta>
                <Card.Description>
                    <MessageHeaderViewer headers={message.Content.Headers} />
                </Card.Description>
                <div style={{ marginTop: "0.5em" }}>
                    <MessageContentViewer message={message}></MessageContentViewer>
                </div>
            </Card.Content>
        </Card>
    );
};

export default MessageViewer;
