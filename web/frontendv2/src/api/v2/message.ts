import axios from "axios";
import { Message } from "../domain";

export interface PageMessageResponse {
    total: number;
    count: number;
    start: number;
    items: Message[];
}

export const getMessages = async (start: number, limit: number): Promise<PageMessageResponse> =>
    axios.get("/api/v2/messages", {
        params: { start, limit },
    });
