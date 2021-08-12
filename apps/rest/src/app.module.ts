import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ConfigModule } from '@nestjs/config';
import CustomTypeOrmModule from './customTypeOrmModule';

@Module({
  imports: [
    CustomTypeOrmModule.forRootAsync(null, {
      entitiesFunc: () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return require.context('./', true, /\.model\.ts/);
      },
      migrationsFunc: () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return require.context('./', true, /\.migration\.ts/);
      },
    }),
    ConfigModule.forRoot({ isGlobal: true }),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
