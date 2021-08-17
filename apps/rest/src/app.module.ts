import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppService } from './services/app';
import CustomTypeOrmModule from './providers/typeorm.module';
import entities from './models';
import { AppController } from './controllers/app';
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
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
