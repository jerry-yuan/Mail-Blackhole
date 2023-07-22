import { Button, Container, Header, Segment } from "semantic-ui-react";
import { Link } from "react-router-dom";
export default function NotFoundPage() {
    return (
        <Segment inverted vertical textAlign="center" style={{ height: "100%" }}>
            <Container>
                <Header as="h1" inverted>
                    404 Not Found
                </Header>
                <p>Oops! Looks like the page you are looking for either does not exist or has been moved.</p>
                <Button primary as={Link} to="/">
                    Go Home
                </Button>
            </Container>
        </Segment>
    );
}
