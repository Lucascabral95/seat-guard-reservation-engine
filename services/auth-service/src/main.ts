import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { envs } from 'config/envs.schemas';
import corsOptions from 'config/cors';
import { ValidationPipe } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.enableCors(corsOptions);
  
  app.useGlobalPipes(
  new ValidationPipe({
 whitelist: true,
 forbidNonWhitelisted: true,
 transform: true,
 })
);

console.log(`Listening on port ${envs.port ?? 3000}`);
await app.listen(envs.port ?? 3000, "0.0.0.0");
}
bootstrap();
