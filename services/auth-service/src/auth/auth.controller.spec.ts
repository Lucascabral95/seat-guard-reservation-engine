import { Test, TestingModule } from '@nestjs/testing';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { CreateAuthDto, LoginAuthDto } from './dto';
import { JwtService } from '@nestjs/jwt'; 
import { AuthGuard } from './auth.guard';

describe('AuthController', () => {
  let controller: AuthController;
  let service: AuthService;

  const mockAuthService = {
    register: jest.fn(),
    login: jest.fn(),
    findAll: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [AuthController],
      providers: [
        {
          provide: AuthService,
          useValue: mockAuthService,
        },
        {
          provide: JwtService,
          useValue: {
            verifyAsync: jest.fn(),
            signAsync: jest.fn(),
          },
        },
      ],
    })
    .overrideGuard(AuthGuard)
    .useValue({ canActivate: () => true }) 
    .compile();

    controller = module.get<AuthController>(AuthController);
    service = module.get<AuthService>(AuthService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('register', () => {
    it('debería llamar a authService.register y retornar el usuario creado', async () => {
      const dto: CreateAuthDto = { 
        email: 'test@test.com', 
        password: 'Password123!', 
        name: 'Test User' 
      };
      
      const expectedResult = { 
        id: 1, 
        email: 'test@test.com', 
        name: 'Test User', 
        createdAt: new Date() 
      } as any;

      mockAuthService.register.mockResolvedValue(expectedResult);

      const result = await controller.register(dto);

      expect(result).toEqual(expectedResult);
      expect(service.register).toHaveBeenCalledWith(dto);
    });
  });

  describe('login', () => {
    it('debería llamar a authService.login y retornar el token', async () => {
      const dto: LoginAuthDto = { email: 'test@test.com', password: 'Password123!' };
      const expectedResult = { access_token: 'fake_jwt_token' };

      mockAuthService.login.mockResolvedValue(expectedResult);

      const result = await controller.login(dto);

      expect(result).toEqual(expectedResult);
      expect(service.login).toHaveBeenCalledWith(dto);
    });
  });
});