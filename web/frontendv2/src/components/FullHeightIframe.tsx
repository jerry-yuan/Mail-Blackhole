import React, { Component, RefObject } from "react";
import ReactDOM from "react-dom";

interface FullheightIframeProps {
    srcDoc: string;
}

interface FullheightIframeState {
    iFrameHeight: string;
}

class FullHeightIframe extends Component<FullheightIframeProps, FullheightIframeState> {
    private iframeRef: RefObject<HTMLIFrameElement>;

    constructor(props: FullheightIframeProps) {
        super(props);
        this.state = {
            iFrameHeight: "0px",
        };
        this.iframeRef = React.createRef();
    }

    componentDidMount() {
        const iframe = this.iframeRef.current;
        if (iframe?.contentWindow?.document?.body) {
            const computedStyle = window.getComputedStyle(iframe.contentWindow.document.body);
            var height = parseInt(computedStyle.height) + parseInt(computedStyle.marginTop) + parseInt(computedStyle.marginBottom);
            iframe.style.height = `${height}px`;
        }
    }

    handleLoad = () => {
        const iframe = this.iframeRef.current;
        if (iframe?.contentWindow?.document?.body) {
            const computedStyle = window.getComputedStyle(iframe.contentWindow.document.body);
            var height = parseInt(computedStyle.height) + parseInt(computedStyle.marginTop) + parseInt(computedStyle.marginBottom);
            iframe.style.height = `${height}px`;
        }
    };

    render() {
        return (
            <iframe
                style={{
                    width: "100%",
                    height: this.state.iFrameHeight,
                    border: "0",
                    overflowY: "hidden",
                }}
                target-blank=""
                seamless
                onLoad={this.handleLoad}
                ref={this.iframeRef}
                srcDoc={this.props.srcDoc}
                height={this.state.iFrameHeight}
            />
        );
    }
}

export default FullHeightIframe;
