import {
  ArgumentsHost,
  Catch,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { EntityNotFoundError, TypeORMError } from 'typeorm';
import { Response } from 'express';

@Catch()
export class GlobalExceptions {
  private logger = new Logger(GlobalExceptions.name);

  catch(exception: Error, host: ArgumentsHost) {
    // log error
    this.logger.error(exception.message, exception.stack, exception);
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();

    if (exception instanceof HttpException) {
      response.status(exception.getStatus()).json(exception.getResponse());
    } else if (exception instanceof TypeORMError) {
      this.handleTypeOrmError(exception, response);
    }
    // add other error handlers here
    else {
      response.status(HttpStatus.INTERNAL_SERVER_ERROR).json({
        message: exception.message,
      });
    }
  }

  handleTypeOrmError(exception: TypeORMError, response: Response) {
    if (exception instanceof EntityNotFoundError) {
      response.status(HttpStatus.NOT_FOUND).json({
        message: exception.message,
      });
    }
    // add more typeorm exceptions here
  }
}
