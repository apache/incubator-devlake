import { Module, ValidationPipe } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppService } from './services/app';
import CustomTypeOrmModule from './providers/typeorm.module';
import entities from './models';
import { AppController } from './controllers/app';
import { APP_FILTER, APP_PIPE } from '@nestjs/core';
import { GlobalExceptions } from './providers/exception';
import { SourceController } from './controllers/source';
import { SourceTaskController } from './controllers/sourceTask';
import { SourceService } from './services/source';
import { SourceTaskService } from './services/sourceTask';
@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: process.env.NODE_ENV === 'test' ? '.env.test' : '.env',
    }),
    CustomTypeOrmModule.forRootAsync(null, {
      entities,
      // FIXME: using db migration instead of synchronize
      synchronize: true,
    }),
  ],
  controllers: [AppController, SourceController, SourceTaskController],
  providers: [
    { provide: APP_FILTER, useClass: GlobalExceptions },
    {
      provide: APP_PIPE,
      useFactory: () => {
        return new ValidationPipe({ transform: true });
      },
    },
    AppService,
    SourceService,
    SourceTaskService,
  ],
})
export class AppModule {}
