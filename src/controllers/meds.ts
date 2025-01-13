import { Hono } from "hono";
import {  readAndSaveFile, saveMedsToDb, searchByName } from "../modules/meds.js";

const meds = new Hono()

meds.get('/', async (c) => {
  const key = c.req.query('key')

  if (key !== process.env.KEY) {
    return c.json({ message: 'Invalid key' }, 401)
  }

  await readAndSaveFile()
  return c.json({ message: 'Meds saved' })
})

meds.get('/search', async (c) => {
  const name = c.req.query('q')
  if (!name) {
    return c.json({ message: 'Query param "q" is required' }, 400)
    
  }
  const meds = await searchByName(name)
  return c.json(meds)
})

meds.post('/', async (c) => {
    return c.json({ message: 'Hello from meds!' })
})

export { meds }