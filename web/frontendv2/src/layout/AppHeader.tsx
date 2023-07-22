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
                </Container>
            </Menu>
        </div>
    );
}
