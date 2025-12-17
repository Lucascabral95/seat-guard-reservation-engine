import { ApiProperty } from "@nestjs/swagger";
import { IsDate, IsString } from "class-validator";

export class GetAuthDto {
    @ApiProperty({
        example: "27abaea9-d2f8-2222-a9b4-a9595516516a",
        description: "User id with uuid format",
        required: false
    })
    @IsString()
   id: string;
   
   @ApiProperty({
        example: "lucas@hotmail.com",
        description: "User email",
        required: true
    })
   @IsString()
   email: string;
   
   @ApiProperty({
        example: "Lucas Doe",
        description: "User name",
        required: false
    })
   @IsString()
   name?: string;

   @ApiProperty({
        example: "2022-01-01T00:00:00.000Z",
        description: "User created at",
        required: true
    })
   @IsDate()
   createdAt: Date;
}