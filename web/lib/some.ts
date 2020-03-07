interface ISome<T, E> {
  Success?: T;
  Error?: E;
}

export type Whether<T, E> = {
  success: (fn: (e: T) => void) => OnError<E>;
};

export type OnError<E> = {
  error: (fn: (e: E) => void) => void;
};

export type Some<T, E> = (data: Whether<T, E>) => OnError<E>;

function success<T, E>(data: ISome<T, E>): Whether<T, E> {
  return {
    success: fn => {
      if (data.Error === null && data.Success != null) {
        fn(data.Success);
        return {
          error: () => {}
        };
      }
      return {
        error: fn => {
          fn(data.Error);
        }
      };
    }
  };
}

export const whether: <T, E>(s: T, e: E) => Whether<T, E> = (s, e) => {
  return success({ Success: s, Error: e });
};
