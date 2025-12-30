import "dotenv/config"
import * as joi from "joi"
import type { StringValue } from "ms";

interface EnvInterface {
  PORT: number
  DATABASE_URL: string
  MY_FRONTEND_URL: string
  SECRET_JWT: string
  EXPIRED_TOKEN_JWT: StringValue
  NODE_ENV: string
}

const varsSchema = joi.object<EnvInterface>({
  PORT: joi.number().required(),
  DATABASE_URL: joi.string().required(),
  MY_FRONTEND_URL: joi.string().required(),
  SECRET_JWT: joi.string().required(),
  EXPIRED_TOKEN_JWT: joi.string().required(),
  NODE_ENV: joi.string().valid('development', 'production', 'test').default('development')
}).unknown(true)

const { error, value: vars } = varsSchema.validate(process.env)

if (error) {
   throw new Error("Invalid environment variables")
}

export const envs = {
  port: vars.PORT,
  databaseUrl: vars.DATABASE_URL,
  myFrontendUrl: vars.MY_FRONTEND_URL,
  secretJwt: vars.SECRET_JWT,
  expiredTokenJwt: vars.EXPIRED_TOKEN_JWT as StringValue,
  nodeEnv: vars.NODE_ENV,
}