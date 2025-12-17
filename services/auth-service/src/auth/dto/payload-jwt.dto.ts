import { ApiProperty } from "@nestjs/swagger";
import { IsDate, IsString } from "class-validator";

export class PayloadJWTDto {
    @ApiProperty({
        example: "27abaea9-d2f8-2222-a9b4-a9595516516a",
        description: "User id with uuid format",
        required: false
    })
    @IsString()
    id: string;

    @ApiProperty({
        example: "lucas@hotmail.com",
        description: "User email with uuid format",
        required: true 
    })
    @IsString()
    email: string;

    @IsDate()
    createdAt: Date;
}