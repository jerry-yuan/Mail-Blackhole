import { Container, Segment, Icon } from "semantic-ui-react";

export function AppFooter() {
    return (
        <div className="app-footer">
            <Segment inverted vertical>
                <Container>
                    <Icon name="copyright outline"></Icon> Copyright 2023 Mailblackhole
                </Container>
            </Segment>
        </div>
    );
}
