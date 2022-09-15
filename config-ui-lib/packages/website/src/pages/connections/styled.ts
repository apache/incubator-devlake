import styled from '@emotion/styled';

export const Container = styled.div`
  h4 {
    margin-top: 12px;
  }

  .item {
    margin-top: 48px;
  }

  ul.list {
    margin-top: 24px;

    li {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 4px 6px;
      width: 130px;
      text-align: center;
      transition: all 0.3s ease;
      cursor: pointer;

      &:hover {
        box-shadow: 1px 1px 6px rgb(0 0 0 / 10%);
      }
    }
  }
`;
