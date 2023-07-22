import axios from "axios";
import { Message } from "../domain";
export const getMessage = async (id: string): Promise<Message> => axios.get("/api/v1/messages/" + id);
export const deleteMessages = async (): Promise<null> => axios.delete("/api/v1/messages");
