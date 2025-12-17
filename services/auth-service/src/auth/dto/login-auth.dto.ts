import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class LoginAuthDto {
    @ApiProperty({
        example: "lucas@hotmail.com",
        description: "User email",
        required: true
    })
    @IsString()
    email: string;
    
    @ApiProperty({
        example: "User.dlkfmlsmdf23",
        description: "User password",
        required: true
    })
    @IsString()
    password: string;   
}