import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import CustomTypeOrmModule from './providers/typeorm.module';
import entities from './models';
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
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
