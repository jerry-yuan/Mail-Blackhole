import React from "react";
import { Link } from "react-router-dom";
import { Button, Container, Header, Icon, Message, Segment } from "semantic-ui-react";
interface PropsType {
    children: React.ReactNode;
}
interface StateType {
    hasError: boolean;
    error?: null | Error;
    errorInfo?: null | React.ErrorInfo;
}
export class ErrorBoundary extends React.Component<PropsType, StateType> {
    constructor(props: PropsType) {
        super(props);
        this.state = {
            hasError: false,
            error: null,
            errorInfo: null,
        };
    }

    static getDerivedStateFromError(error: Error): StateType {
        return { hasError: true };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        this.setState((preState) => ({ hasError: preState.hasError, error: error, errorInfo: errorInfo }));
    }
    render() {
        if (this.state.hasError) {
            return (
                <Segment inverted vertical style={{ height: "100%" }}>
                    <Container textAlign="center">
                        <Header as="h1" inverted>
                            Oops! Page Crashed!
                        </Header>
                        <p>There are some errors when rendering the page you are looking for.</p>
                        <Button primary as={Link} to="/">
                            Refresh
                        </Button>
                    </Container>
                    <Container style={{ marginTop: "1em" }}>
                        <Message icon="bug" color="yellow" header={this.state.error}>
                            <pre>
                                {this.state.error && this.state.error.message}
                                {this.state.errorInfo?.componentStack}
                            </pre>
                        </Message>
                    </Container>
                </Segment>
            );
        }
        return this.props.children;
    }
}
