import { Module, ValidationPipe } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppService } from './services/app';
import CustomTypeOrmModule from './providers/typeorm.module';
import entities from './models';
import { AppController } from './controllers/app';
import { APP_FILTER, APP_PIPE } from '@nestjs/core';
import { NotFoundFilter } from './providers/exception';
import { SourceController } from './controllers/source';
import { SourceService } from './services/source';
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
  controllers: [AppController, SourceController],
  providers: [
    { provide: APP_FILTER, useClass: NotFoundFilter },
    {
      provide: APP_PIPE,
      useFactory: () => {
        return new ValidationPipe({ transform: true });
      },
    },
    AppService,
    SourceService,
  ],
})
export class AppModule {}
