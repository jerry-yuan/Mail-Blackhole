import { Container, Menu, Image, Dropdown } from "semantic-ui-react";
import logo from "../blackhole.svg";
import { Link } from "react-router-dom";
export default function AppHeader() {
    return (
        <div className="app-header">
            <Menu inverted style={{ borderRadius: 0 }}>
                <Container>
                    <Menu.Item as={Link} to="/" header>
                        <Image size="mini" src={logo} style={{ marginRight: "1.5em" }} />
                        MailBlackhole
                    </Menu.Item>
                    <Menu.Item as={Link} to="/">
                        Home
                    </Menu.Item>

                    <Dropdown item simple text="Dropdown">
                        <Dropdown.Menu>
                            <Dropdown.Item>List Item</Dropdown.Item>
                            <Dropdown.Item>List Item</Dropdown.Item>
                            <Dropdown.Divider />
                            <Dropdown.Header>Header Item</Dropdown.Header>
                            <Dropdown.Item>
                                <i className="dropdown icon" />
                                <span className="text">Submenu</span>
                                <Dropdown.Menu>
                                    <Dropdown.Item>List Item</Dropdown.Item>
                                    <Dropdown.Item>List Item</Dropdown.Item>
                                </Dropdown.Menu>
                            </Dropdown.Item>
                            <Dropdown.Item>List Item</Dropdown.Item>
                        </Dropdown.Menu>
                    </Dropdown>
                </Container>
            </Menu>
        </div>
    );
}
