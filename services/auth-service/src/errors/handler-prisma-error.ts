import { NotFoundException, BadRequestException, InternalServerErrorException } from '@nestjs/common';

export function handlePrismaError(error: any, entityName: string, id?: string) {
    switch (error.code) {
        case 'P2025':
            {
                const message = id
                    ? `${entityName} with id: ${id} not found`
                    : `${entityName} not found`;
                throw new NotFoundException(message);
            }

        case 'P2002':
            throw new BadRequestException(`${entityName} already exists`);

        case 'P2003':
            throw new BadRequestException(`Cannot delete ${entityName}: has related records`);

        default:
            throw new InternalServerErrorException(`Error processing ${entityName}`);
    }
}