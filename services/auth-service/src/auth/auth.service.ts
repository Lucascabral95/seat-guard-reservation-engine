import { BadRequestException, Injectable, InternalServerErrorException, NotFoundException, UnauthorizedException } from '@nestjs/common';
import { CreateAuthDto, LoginAuthDto, PayloadJWTDto } from './dto';
import * as bcrypt from "bcrypt"
import { PrismaService } from 'src/prisma/prisma.service';
import { handlePrismaError } from 'src/errors/handler-prisma-error';
import { JwtService } from '@nestjs/jwt';
import { COMMON_PASSWORDS } from 'src/common/constants/security.constants';

@Injectable()
export class AuthService {

  constructor(
    private readonly prisma: PrismaService,
    private jwtService: JwtService
  ) {}

 async register(createAuthDto: CreateAuthDto) {
    try {
      const { email, password, name } = createAuthDto;

      await this.findUserByEmail(email);

      const normalizedPassword = password.toLowerCase();

      if (COMMON_PASSWORDS.has(normalizedPassword)) {
        throw new BadRequestException('Esta contraseña es demasiado común');
      }

      const saltOrRounds = 12;
      const passwordHashed = await bcrypt.hash(password, saltOrRounds);
      const lowerCaseEmail = email.toLowerCase();

      const user = await this.prisma.user.create({
        data: {
          email: lowerCaseEmail,
          password: passwordHashed,
          name,
        },
      });

      const { password: _, ...userWithoutPassword } = user;
      return userWithoutPassword;
    } catch (error) {
      if (
        error instanceof BadRequestException ||
        error instanceof InternalServerErrorException ||
        error instanceof NotFoundException
      ) {
        throw error;
      }
      handlePrismaError(error, 'User');
    }
  }

   async login(loginAuthDto: LoginAuthDto) {
  try {
    const { email, password } = loginAuthDto;
    const lowerCaseEmail = email.toLowerCase()

    const user = await this.prisma.user.findFirst({
      where: { 
        email: lowerCaseEmail,
       },
    });

    if (!user) {
      throw new UnauthorizedException('Credenciales inválidas.');
    }

    const isPasswordValid = await bcrypt.compare(password, user.password);
    if (!isPasswordValid) {
      throw new UnauthorizedException('Credenciales inválidas.');
    }

    const res = this.signIn({
      id: user.id,
      email: user.email,
      createdAt: user.createdAt,
    })

    return res;
  } catch (error) {
    if (
      error instanceof UnauthorizedException ||
      error instanceof BadRequestException ||
      error instanceof NotFoundException ||
      error instanceof InternalServerErrorException
    ) {
      throw error;
    }
    handlePrismaError(error, 'User');
  }
}

async signIn(payloaJWTDto: PayloadJWTDto) {
  const user = payloaJWTDto

  if (!user) {
    throw new UnauthorizedException();
  }
  
  return {
    access_token: await this.jwtService.signAsync(payloaJWTDto),
  }
}

  async findAll() {
    try {
      const users = await this.prisma.user.findMany();

      if (!users) {
        return new NotFoundException("Sin usuarios registrados.")
      }

      return users;
    } catch (error) {
       if (error instanceof UnauthorizedException || error instanceof BadRequestException || error instanceof NotFoundException || error instanceof InternalServerErrorException ) {
        throw error;
       }
       handlePrismaError(error, "Users")
    }
  }
  
  async findUserByEmail(email: string) {
    try {
       const findUser = await this.prisma.user.findFirst({
        where: {
          email: email,
        },
       })

       if (findUser) {
        throw new BadRequestException("Usuario ya existente con este email")
       }
      
       return;
    } catch (error) {
      if(error instanceof UnauthorizedException || error instanceof BadRequestException || error instanceof NotFoundException || error instanceof InternalServerErrorException ) {
        throw error;
      }
      handlePrismaError(error, "Users")
    }
  }

}