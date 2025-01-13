
import fs from 'fs';
import { db } from '../db/config.js';
import { meds } from '../db/schema.js';
import {  createInterface } from 'node:readline/promises';
import {like, or } from 'drizzle-orm';


const CSV_PATH = './assets/DADOS_ABERTOS_MEDICAMENTOS(in).csv';

export async function readAndSaveFile() {
  //read file line by line
  const fileStream = fs.createReadStream(CSV_PATH);
  const rl = createInterface({
    input: fileStream,
    crlfDelay: Infinity
  })

  let isFirstLine = true;
  let headerArray: string[] = [];


  for await (const line of rl) {
    if (isFirstLine) {
      headerArray = line.split(';').map(header => header.trim());
      isFirstLine = false;
      continue;
    }

    const rowData = line.split(';');
    const formattedData = headerArray.reduce((acc: Record<string, string>, header, i) => {
      acc[header] = rowData[i]?.trim() || '';
      return acc;
    }, {});

    await saveMedsToDb([formattedData]);
  }

  rl.close();



}

export async function saveMedsToDb(data: Record<string, string>[]) {
  return await db.insert(meds).values(data)
}

export async function searchByName(name: string) {
  const result = await db
    .selectDistinct(({ NOME_PRODUTO: meds.NOME_PRODUTO, PRINCPIO_ATIVO: meds.PRINCIPIO_ATIVO }))
    .from(meds)
    .where(
      or(
        like(meds.PRINCIPIO_ATIVO, `%${name}%`),
        like(meds.NOME_PRODUTO, `%${name}%`)
      )
    )
    .all();

  return result;



}