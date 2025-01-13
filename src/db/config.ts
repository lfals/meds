import 'dotenv/config';
import { drizzle } from 'drizzle-orm/better-sqlite3';

export const db = drizzle({logger: true, connection: { source: process.env.DB_FILE_NAME }});