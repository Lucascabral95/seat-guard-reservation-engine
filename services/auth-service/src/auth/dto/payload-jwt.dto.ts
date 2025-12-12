import { IsDate, IsString } from "class-validator";

export class PayloadJWTDto {
    @IsString()
    id: string;

    @IsString()
    email: string;
    
    @IsDate()
    createdAt: Date;
}