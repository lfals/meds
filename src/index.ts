import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { meds } from './controllers/meds'

const app = new Hono()


app.route('/meds', meds)

app.get('/', async (c) => {
  return c.text('Hello Hono!')
})

const port = 3000
console.log(`Server is running on http://localhost:${port}`)

serve({
  fetch: app.fetch,
  port
})
