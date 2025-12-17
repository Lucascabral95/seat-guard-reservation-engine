import { Test, TestingModule } from '@nestjs/testing';
import { JwtService } from '@nestjs/jwt';
import { BadRequestException, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcrypt';
import { AuthService } from './auth.service';
import { PrismaService } from 'src/prisma/prisma.service';

jest.mock('bcrypt', () => ({
  hash: jest.fn(),
  compare: jest.fn(),
}));

describe('AuthService', () => {
  let service: AuthService;
  let prisma: PrismaService;
  let jwtService: JwtService;

  const mockPrismaService = {
    user: {
      create: jest.fn(),
      findFirst: jest.fn(),
      findMany: jest.fn(),
    },
  };

  const mockJwtService = {
    signAsync: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
        {
          provide: JwtService,
          useValue: mockJwtService,
        },
      ],
    }).compile();

    service = module.get<AuthService>(AuthService);
    prisma = module.get<PrismaService>(PrismaService);
    jwtService = module.get<JwtService>(JwtService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('register', () => {
    it('debería registrar un usuario exitosamente', async () => {
      const dto = { email: 'test@test.com', password: 'StrongPassword123!', name: 'Test' };
      const hashedPassword = 'hashed_password';
      
      mockPrismaService.user.findFirst.mockResolvedValue(null);
      (bcrypt.hash as jest.Mock).mockResolvedValue(hashedPassword);
      mockPrismaService.user.create.mockResolvedValue({
        id: 1,
        ...dto,
        password: hashedPassword,
        createdAt: new Date(),
      });

      const result = await service.register(dto);

      expect(mockPrismaService.user.findFirst).toHaveBeenCalledWith({ where: { email: dto.email } });
      expect(bcrypt.hash).toHaveBeenCalledWith(dto.password, 12);
      expect(mockPrismaService.user.create).toHaveBeenCalled();
      expect(result).not.toHaveProperty('password'); 
      expect(result?.email).toBe(dto.email.toLowerCase());
    });

    it('debería lanzar BadRequestException si el email ya existe', async () => {
      const dto = { email: 'exists@test.com', password: 'Password123!', name: 'Test' };
      
      mockPrismaService.user.findFirst.mockResolvedValue({ id: 1, email: dto.email });

      await expect(service.register(dto)).rejects.toThrow(BadRequestException);
    });

    it('debería lanzar BadRequestException si la contraseña es común', async () => {
      const dto = { email: 'new@test.com', password: 'password', name: 'Test' }; 
      
      mockPrismaService.user.findFirst.mockResolvedValue(null);

      await expect(service.register(dto)).rejects.toThrow(BadRequestException);
    });
  });

  describe('login', () => {
    it('debería retornar un token si las credenciales son válidas', async () => {
      const dto = { email: 'test@test.com', password: 'Password123!' };
      const userInDb = { 
        id: 1, 
        email: 'test@test.com', 
        password: 'hashed_password', 
        createdAt: new Date() 
      };

      mockPrismaService.user.findFirst.mockResolvedValue(userInDb);
      (bcrypt.compare as jest.Mock).mockResolvedValue(true); 
      mockJwtService.signAsync.mockResolvedValue('fake_jwt_token');

      const result = await service.login(dto);

      expect(result).toEqual({ access_token: 'fake_jwt_token' });
      expect(mockJwtService.signAsync).toHaveBeenCalled();
    });

    it('debería lanzar UnauthorizedException si el usuario no existe', async () => {
      const dto = { email: 'wrong@test.com', password: 'Password123!' };
      mockPrismaService.user.findFirst.mockResolvedValue(null);

      await expect(service.login(dto)).rejects.toThrow(UnauthorizedException);
    });

    it('debería lanzar UnauthorizedException si la contraseña es incorrecta', async () => {
      const dto = { email: 'test@test.com', password: 'WrongPassword!' };
      const userInDb = { id: 1, email: 'test@test.com', password: 'hashed_password' };

      mockPrismaService.user.findFirst.mockResolvedValue(userInDb);
      (bcrypt.compare as jest.Mock).mockResolvedValue(false); 

      await expect(service.login(dto)).rejects.toThrow(UnauthorizedException);
    });
  });
});