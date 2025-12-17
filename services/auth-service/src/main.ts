import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import corsOptions from 'src/config/cors';
import { ValidationPipe } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { envs } from './config/envs.schemas';

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

const config = new DocumentBuilder()
.setTitle("SeatGuard API")
.setDescription("SeatGuard API description")
.setVersion("1.0")
.addTag("auth")
.addBearerAuth()
.build();

const document = SwaggerModule.createDocument(app, config)
SwaggerModule.setup("api", app, document)

console.log(`Listening on port ${envs.port ?? 3000}`);
await app.listen(envs.port ?? 3000, "0.0.0.0");
}
bootstrap();
