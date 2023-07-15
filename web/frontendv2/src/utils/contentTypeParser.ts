interface ContentType {
    mime: string;
    props: Record<string, string | undefined>;
}

const parse = (contentTypeRaw: string | undefined): ContentType => {
    contentTypeRaw = contentTypeRaw ?? "*/*";
    let contentType = "";
    let contentTypeProps: Record<string, string | undefined> = {};
    let semicolonIndex = contentTypeRaw.indexOf(";");

    // parse Content-Type
    if (semicolonIndex >= 0) {
        contentType = contentTypeRaw.substring(0, semicolonIndex);
        contentTypeRaw = contentTypeRaw.substring(semicolonIndex + 1);
    } else {
        contentType = contentTypeRaw;
        contentTypeRaw = "";
    }
    // parse Content-Type properties
    semicolonIndex = contentTypeRaw.indexOf(";");
    while (contentTypeRaw.length > 0) {
        let nextPropPart;
        if (semicolonIndex >= 0) {
            nextPropPart = contentTypeRaw.substring(0, semicolonIndex);
            contentTypeRaw = contentTypeRaw.substring(semicolonIndex + 1);
        } else {
            nextPropPart = contentTypeRaw;
            contentTypeRaw = "";
        }

        const equalsIndex = nextPropPart.indexOf("=");
        let propName: string, propValue: string | undefined;
        if (equalsIndex < 0) {
            propName = nextPropPart;
            propValue = undefined;
        } else {
            propName = nextPropPart.substring(0, equalsIndex);
            propValue = nextPropPart.substring(equalsIndex + 1);
        }

        contentTypeProps[propName.trim().toLowerCase()] = propValue?.substring(1, propValue.length - 1);
    }
    return {
        mime: contentType,
        props: contentTypeProps,
    };
};

export default { parse };
