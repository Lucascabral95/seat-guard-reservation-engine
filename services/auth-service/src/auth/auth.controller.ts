import { Controller, Get, Post, Body, UseGuards } from '@nestjs/common';
import { AuthService } from './auth.service';
import { CreateAuthDto, LoginAuthDto } from './dto';
import { AuthGuard } from './auth.guard';
import { ApiBadRequestResponse, ApiInternalServerErrorResponse, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags("auth")
@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post("register")
  @ApiOperation({ summary: "Registrar un nuevo usuario" })
  @ApiResponse({ status: 201, description: "Usuario creado exitosamente" })
  @ApiResponse({ status: 400, description: "Error al crear el usuario" })
  @ApiBadRequestResponse({ description: 'Datos inválidos o contraseña débil.' })
  @ApiInternalServerErrorResponse({ description: 'Internal server error' })
  register(@Body() createAuthDto: CreateAuthDto) {
    return this.authService.register(createAuthDto);
  }

  @Post("login")
  @ApiOperation({ summary: "Iniciar sesión" })
  @ApiResponse({ status: 200, description: "Sesión iniciada exitosamente" })
  @ApiResponse({ status: 401, description: "Credenciales inválidas" })
  @ApiInternalServerErrorResponse({ description: 'Internal server error' })
  login(@Body() loginAuthDto: LoginAuthDto) {
    return this.authService.login(loginAuthDto);
  }
  
  @UseGuards(AuthGuard)
  @Get()
  @ApiOperation({ summary: "Obtener todos los usuarios" })
  @ApiResponse({ status: 200, description: "Todos los usuarios obtenidos exitosamente" })
  @ApiInternalServerErrorResponse({ description: 'Internal server error' })
  findAll() {
    return this.authService.findAll();
  }
}
