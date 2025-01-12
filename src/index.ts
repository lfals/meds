import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { db } from './db/config.js'
import { meds } from './controllers/meds.js'

const app = new Hono()


app.route('/meds', meds)

app.get('/', async (c) => {
  const result = await db.run('select 1')
  return c.text('Hello Hono!' + JSON.stringify(result))
})

const port = 3000
console.log(`Server is running on http://localhost:${port}`)

serve({
  fetch: app.fetch,
  port
})
