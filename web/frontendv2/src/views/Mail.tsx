import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Container, Icon, Button, Segment, Dimmer, Loader } from "semantic-ui-react";
import { getMessage } from "../api/v1/messages";
import ErrorPage from "./ErrorPage";
import { Message } from "../api/domain";
import MessageViewer from "../components/message/MessageViewer";

export default function MailPage() {
    const { id } = useParams();
    const [message, setMessage] = useState<Message | null>(null);

    useEffect(() => {
        (async () => {
            if (id) {
                setMessage(await getMessage(id));
            }
        })();
    }, [id]);

    return id ? (
        <Container style={{ marginTop: "1em", marginBottom: "1em" }}>
            <div>
                <Button as={Link} to="/">
                    <Icon name="arrow left"></Icon>
                </Button>
            </div>
            {message !== null ? (
                <MessageViewer message={message}></MessageViewer>
            ) : (
                <Segment>
                    <Dimmer active>
                        <Loader size="large">Loading</Loader>
                    </Dimmer>
                </Segment>
            )}
        </Container>
    ) : (
        <ErrorPage />
    );
}
