export interface Message {
    ID: string;
    From: Path;
    To: Path[];
    Content: Content;
    Created: Date;
    MIME: MIMEBody | null;
    Raw: SMTPMessage | null;
}

export interface Path {
    Relays: string[];
    Mailbox: string;
    Domain: string;
    Params: string;
}

export interface Content {
    Headers: Record<string, string[]>;
    Body: string;
    Size: number;
    MIME: MIMEBody | null;
}

export interface SMTPMessage {
    From: string;
    To: string[];
    Data: string;
    Helo: string;
}

export interface MIMEBody {
    Parts: Content[];
}
