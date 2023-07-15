import { useEffect, useState } from "react";
import { Button, Container, Divider, Dropdown, Header, Icon, List, Pagination, Placeholder, Popup, Segment } from "semantic-ui-react";

import { getMessages } from "../api/v2/message";
import { Message } from "../api/domain";
import { Link } from "react-router-dom";
import moment from "moment";
import prettyBytes from "pretty-bytes";
import Rfc2047Text from "../components/Rfc2047Text";
const availablePageSizes = [
    { text: "7", value: 8 },
    { text: "10", value: 10 },
    { text: "20", value: 20 },
    { text: "50", value: 50 },
    { text: "100", value: 100 },
    { text: "200", value: 200 },
    { text: "500", value: 500 },
];

export default function InboxPage() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [pagination, setPagination] = useState({ totalPages: 0, activePage: 1, pageSize: 8 });
    const [loading, setLoading] = useState(true);

    const refresh = async function () {
        const start = (pagination.activePage - 1) * pagination.pageSize;
        setLoading(true);
        const pageResponse = await getMessages(start, pagination.pageSize);

        setMessages(pageResponse.items);
        setPagination({ ...pagination, totalPages: Math.ceil(pageResponse.total / pagination.pageSize) });
        setLoading(false);
    };

    useEffect(() => {
        refresh();
    }, [pagination.activePage, pagination.pageSize]);
    return (
        <Container style={{ marginTop: "1em", marginBottom: "1em" }}>
            <Header as="h2">Inbox {JSON.stringify(pagination)}</Header>
            <Divider inverted />
            <div style={{ display: "flex", justifyContent: "space-between" }}>
                <div>
                    <Button icon="refresh" color="green" onClick={refresh}></Button>
                    <Button icon="trash" color="red"></Button>
                </div>
                <div>
                    <Dropdown
                        options={availablePageSizes}
                        value={pagination.pageSize}
                        fluid
                        selection
                        onChange={(e, { value }) => setPagination({ ...pagination, pageSize: value as number })}
                    />
                </div>
            </div>
            <Segment>
                {loading ? (
                    <Placeholder fluid>
                        {Array.from({ length: pagination.pageSize }, (t) => (
                            <Placeholder.Paragraph>
                                <Placeholder.Line />
                                <Placeholder.Line />
                            </Placeholder.Paragraph>
                        ))}
                    </Placeholder>
                ) : (
                    <List divided relaxed>
                        {messages.map((message) => (
                            <List.Item key={message.ID}>
                                <List.Icon name="mail" size="large" verticalAlign="middle" />
                                <List.Content>
                                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                                        <div style={{ display: "flex", justifyContent: "left" }}>
                                            <List.Header as={Link} to={"/mail/" + message.ID}>
                                                <Rfc2047Text>{message.Content.Headers["Subject"]}</Rfc2047Text>
                                            </List.Header>
                                            <List.Description></List.Description>
                                        </div>
                                        <List.Description>{moment(message.Created).format("yyyy-MM-DD HH:mm:ss")}</List.Description>
                                    </div>
                                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                                        <div>
                                            <Popup position="bottom center" trigger={<Icon className="purple" name="user circle" />}>
                                                <Rfc2047Text>
                                                    From: {message.From.Mailbox}@{message.From.Domain}
                                                </Rfc2047Text>
                                            </Popup>
                                            {message.To.map((receiver, index) => (
                                                <Popup key={index} position="bottom center" trigger={<Icon className="blue" name="user circle" />}>
                                                    To:{receiver.Mailbox}@{receiver.Domain}
                                                </Popup>
                                            ))}
                                        </div>
                                        <div style={{ display: "flex", justifyContent: "right" }}>{prettyBytes(message.Raw?.Data.length || 0)}</div>
                                    </div>
                                </List.Content>
                            </List.Item>
                        ))}
                    </List>
                )}
            </Segment>
            <div style={{ display: "flex", justifyContent: "right" }}>
                <Pagination
                    totalPages={pagination.totalPages}
                    activePage={pagination.activePage}
                    onPageChange={(e, { activePage }) => setPagination({ ...pagination, activePage: activePage as number })}
                />
            </div>
        </Container>
    );
}
