import { ApiProperty } from "@nestjs/swagger";
import { IsEmail, IsOptional, IsString } from "class-validator";
import { IsStrongPassword } from "src/common/validators/strong-password-validator";

export class CreateAuthDto {
    @ApiProperty({
        example: "lucas@hotmail.com",
        description: "User email",
        required: true
    })
    @IsString()
    @IsEmail()
    email: string;
    
    @ApiProperty({
        example: "User.dlkfmlsmdf23",
        description: "User password",
        required: true
    })
    @IsString()
    @IsStrongPassword()
    password: string;
    
    @ApiProperty({
        example: "Lucas Doe",
        description: "User name",
        required: true
    })
    @IsString()
    @IsOptional()
    name?: string
}
