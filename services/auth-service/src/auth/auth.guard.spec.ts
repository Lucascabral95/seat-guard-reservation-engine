import { Test, TestingModule } from '@nestjs/testing';
import { AuthGuard } from './auth.guard';
import { JwtService } from '@nestjs/jwt';
import { UnauthorizedException } from '@nestjs/common';
import { ExecutionContext } from '@nestjs/common';

jest.mock('src/config/envs.schemas', () => ({
  envs: {
    secretJwt: 'test-secret',
  },
}));

describe('AuthGuard', () => {
  let guard: AuthGuard;
  let jwtService: JwtService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthGuard,
        {
          provide: JwtService,
          useValue: {
            verifyAsync: jest.fn(),
          },
        },
      ],
    }).compile();

    guard = module.get<AuthGuard>(AuthGuard);
    jwtService = module.get<JwtService>(JwtService);
  });

  const createMockContext = (authHeader?: string): ExecutionContext => {
    return {
      switchToHttp: () => ({
        getRequest: () => ({
          headers: {
            authorization: authHeader,
          },
        }),
      }),
    } as any;
  };

  it('debería estar definido', () => {
    expect(guard).toBeDefined();
  });

  it('debería permitir el acceso si el token es válido', async () => {
    const context = createMockContext('Bearer valid_token');
    
    jest.spyOn(jwtService, 'verifyAsync').mockResolvedValue({ sub: 1, email: 'test@test.com' });

    const result = await guard.canActivate(context);

    expect(result).toBe(true);
    expect(jwtService.verifyAsync).toHaveBeenCalledWith('valid_token', { secret: 'test-secret' });
  });

  it('debería lanzar UnauthorizedException si no hay cabecera de autorización', async () => {
    const context = createMockContext(undefined);
    await expect(guard.canActivate(context)).rejects.toThrow(UnauthorizedException);
  });

  it('debería lanzar UnauthorizedException si el formato del token es incorrecto', async () => {
    const context = createMockContext('Basic token_invalido');
    await expect(guard.canActivate(context)).rejects.toThrow(UnauthorizedException);
  });

  it('debería lanzar UnauthorizedException si el token es inválido', async () => {
    const context = createMockContext('Bearer invalid_token');
    jest.spyOn(jwtService, 'verifyAsync').mockRejectedValue(new Error());
    await expect(guard.canActivate(context)).rejects.toThrow(UnauthorizedException);
  });
});