import { int, sqliteTable, text } from "drizzle-orm/sqlite-core"



export const meds = sqliteTable("meds", {
    id: int().primaryKey({ autoIncrement: true }),
    TIPO_PRODUTO: text(),
    NOME_PRODUTO: text(),
    DATA_FINALIZACAO_PROCESSO: text(),
    CATEGORIA_REGULATORIA: text(),
    NUMERO_REGISTRO_PRODUTO: text(),
    DATA_VENCIMENTO_REGISTRO: text(),
    NUMERO_PROCESSO: text(),
    CLASSE_TERAPEUTICA: text(),
    EMPRESA_DETENTORA_REGISTRO: text(),
    SITUACAO_REGISTRO: text(),
    PRINCIPIO_ATIVO: text(),
})

