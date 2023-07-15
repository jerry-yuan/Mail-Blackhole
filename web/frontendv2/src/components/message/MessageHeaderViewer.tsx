import { useState } from "react";
import { Collapse } from "react-collapse";
import { Button, Icon, Label, List } from "semantic-ui-react";
import Rfc2047Text from "../Rfc2047Text";

interface PropsType {
    headers: Record<string, string[]>;
}

const MessageHeaderViewer: React.FC<PropsType> = function ({ headers }) {
    const headerKeys = Object.keys(headers);
    const [collapsed, setCollapsed] = useState<boolean>(true);
    return (
        <>
            <Collapse isOpened={!collapsed}>
                <List divided selection>
                    {headerKeys.map((headerKey, index) => {
                        return (
                            <List.Item key={index}>
                                <Label horizontal>{headerKey}</Label>
                                {headers[headerKey].map((v, index) => (
                                    <Rfc2047Text key={index}>{v}</Rfc2047Text>
                                ))}
                            </List.Item>
                        );
                    })}
                </List>
            </Collapse>
            <Button size="tiny" className="fluid" onClick={() => setCollapsed(!collapsed)}>
                <Icon name={collapsed ? "chevron down" : "chevron up"} />
                {collapsed ? "Show All Headers" : "Collapse Headers"}
            </Button>
        </>
    );
};

export default MessageHeaderViewer;
