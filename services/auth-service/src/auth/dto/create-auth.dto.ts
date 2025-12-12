import { IsEmail, IsOptional, IsString, MinLength } from "class-validator";

export class CreateAuthDto {
    @IsString()
    @IsEmail()
    email: string;
    
    @IsString()
    @MinLength(6)
    password: string;
    
    @IsString()
    @IsOptional()
    name?: string
}
